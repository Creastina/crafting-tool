using Microsoft.AspNetCore.Mvc;
using Microsoft.AspNetCore.Mvc.Filters;

namespace Creastina.CraftingTool.Authentication;

public class ApiKeyFilter(IConfiguration configuration) : IAuthorizationFilter
{
    public void OnAuthorization(AuthorizationFilterContext context)
    {
        var configuredApiKey = configuration["ApiKey"];
        var authorizationHeader =
            context.HttpContext.Request.Headers.Authorization.FirstOrDefault()?.Remove(0, "Bearer ".Length);

        if (authorizationHeader == null || configuredApiKey != authorizationHeader)
            context.Result = new UnauthorizedResult();
    }
}