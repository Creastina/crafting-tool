namespace Creastina.CraftingTool.Authentication;

public class CookieAuthenticationChecker(IHttpContextAccessor httpContextAccessor, IConfiguration configuration)
{
    public void CheckCookie()
    {
        if (httpContextAccessor.HttpContext?.Request.Cookies["ApiKey"] != configuration["ApiKey"])
        {
            httpContextAccessor.HttpContext?.Response.Redirect("/login");
        }
    }
}