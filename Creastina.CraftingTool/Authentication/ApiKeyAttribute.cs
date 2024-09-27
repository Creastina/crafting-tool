using Microsoft.AspNetCore.Mvc;

namespace Creastina.CraftingTool.Authentication;

public class ApiKeyAttribute() : ServiceFilterAttribute(typeof(ApiKeyFilter));