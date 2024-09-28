using Microsoft.AspNetCore.Mvc;
using Microsoft.AspNetCore.Mvc.RazorPages;

namespace Creastina.CraftingTool.Pages;

public class Login(IHttpContextAccessor httpContextAccessor, IConfiguration configuration) : PageModel
{
    [BindProperty] public string Password { get; set; }

    public IActionResult OnPost()
    {
        if (configuration["ApiKey"] != Password) return Unauthorized();

        httpContextAccessor.HttpContext.Response.Cookies.Append("ApiKey", Password, new CookieOptions
        {
            Expires = DateTime.Now.AddYears(100)
        });
        return Redirect("/");
    }
}