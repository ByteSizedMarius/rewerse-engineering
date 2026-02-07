"""Low-level cffi bindings to the rewerse Go library."""

import json
import sys
from pathlib import Path

from cffi import FFI

ffi = FFI()

ffi.cdef("""
    void FreeString(char* s);

    // Init
    char* SetCertificate(char* certPath, char* keyPath);

    // Markets
    char* MarketSearch(char* query);
    char* GetMarketDetails(char* marketID);

    // Products
    char* GetProducts(char* marketID, char* search, char* optsJSON);
    char* GetCategoryProducts(char* marketID, char* categorySlug, char* optsJSON);
    char* GetProductByID(char* marketID, char* productID);
    char* GetProductSuggestions(char* query, char* optsJSON);
    char* GetProductRecommendations(char* marketID, char* listingID);

    // Discounts
    char* GetDiscountsRaw(char* marketID);
    char* GetDiscounts(char* marketID);

    // Recipes
    char* RecipeSearch(char* optsJSON);
    char* GetRecipeDetails(char* recipeID);
    char* GetRecipePopularTerms();

    // Misc
    char* GetRecalls();
    char* GetRecipeHub();
    char* GetServicePortfolio(char* zipcode);
    char* GetShopOverview(char* marketID);
    char* GetShopOverviewWithOpts(char* marketID, char* optsJSON);

    // Basket
    char* CreateBasket(char* marketID, char* zipCode, char* serviceType);
    char* GetBasket(char* basketID, char* marketID, char* zipCode, char* serviceType, int version);
    char* SetBasketItemQuantity(char* basketID, char* marketID, char* zipCode, char* serviceType, char* listingID, int quantity, int version);
    char* RemoveBasketItem(char* basketID, char* marketID, char* zipCode, char* serviceType, char* listingID);

    // Delivery
    char* GetBulkyGoodsConfig(char* marketID, char* serviceType);
""")


def _find_library() -> Path:
    """Find the shared library based on platform."""
    lib_dir = Path(__file__).parent / "_lib"

    if sys.platform == "win32":
        lib_name = "rewerse.dll"
    elif sys.platform == "darwin":
        lib_name = "librewerse.dylib"
    else:
        lib_name = "librewerse.so"

    lib_path = lib_dir / lib_name
    if not lib_path.exists():
        raise RuntimeError(
            f"Shared library not found at {lib_path}. "
            f"Run the build script first: python build_lib.py"
        )

    return lib_path


_lib = None

# RTLD_NODELETE prevents the library from being unloaded on dlclose.
# Go's runtime doesn't support being unloaded, so this prevents crashes on exit.
# Linux: Go linker adds -Wl,-z,nodelete automatically
# macOS: Must specify RTLD_NODELETE when loading
# Windows: Handled via DLL pinning in Go's init()
_RTLD_NODELETE = 0x80 if sys.platform == "darwin" else 0


def get_lib():
    """Get the loaded library (lazy loading)."""
    global _lib
    if _lib is None:
        lib_path = _find_library()
        flags = ffi.RTLD_NOW | _RTLD_NODELETE
        _lib = ffi.dlopen(str(lib_path), flags)
    return _lib


class RewerseError(Exception):
    """Error from the rewerse library."""
    pass


def call(func_name: str, *args) -> dict:
    """
    Call a library function and parse the JSON result.

    Args:
        func_name: Name of the C function to call
        *args: Arguments to pass (strings encoded to bytes, ints passed directly)

    Returns:
        The parsed 'data' field from the response

    Raises:
        RewerseError: If the library returns an error
    """
    lib = get_lib()
    func = getattr(lib, func_name)

    # Convert args to C types
    # Note: cffi manages c_args memory - they're freed when c_args goes out of scope
    c_args = []
    for arg in args:
        if arg is None:
            c_args.append(ffi.NULL)
        elif isinstance(arg, int):
            c_args.append(arg)
        else:
            c_args.append(ffi.new("char[]", str(arg).encode("utf-8")))

    result_ptr = func(*c_args)

    # Parse and free the result, preserving original exception if any
    original_error = None
    result = None
    try:
        result_str = ffi.string(result_ptr).decode("utf-8")
        result = json.loads(result_str)
    except Exception as e:
        original_error = e
    finally:
        try:
            lib.FreeString(result_ptr)
        except Exception:
            pass  # Suppress cleanup errors to preserve original exception

    if original_error is not None:
        raise original_error

    if not result.get("ok"):
        raise RewerseError(result.get("error") or "Unknown error")

    return result.get("data")
