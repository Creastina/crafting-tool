using Creastina.CraftingTool.Authentication;
using Creastina.CraftingTool.Models;
using Creastina.CraftingTool.Repository;
using Creastina.CraftingTool.Components;
using Microsoft.EntityFrameworkCore;

var builder = WebApplication.CreateBuilder(args);

// Add services to the container.
builder.Services.AddControllersWithViews();

builder.Services.AddEndpointsApiExplorer();
builder.Services.AddSwaggerGen();
builder.Configuration.AddEnvironmentVariables();

builder.Services.AddDbContext<CraftingContext>(optionsBuilder =>
    optionsBuilder.UseNpgsql(builder.Configuration.GetConnectionString("Crafting")));
builder.Services.AddNpgsql<CraftingContext>(builder.Configuration.GetConnectionString("Crafting"));

builder.Services.AddScoped<ApiKeyFilter>(); 
builder.Services.AddScoped<CookieAuthenticationChecker>();
builder.Services.AddScoped<ITodoRepository, TodoRepository>();

builder.Services.AddHttpContextAccessor();
builder.Services.AddRazorComponents()
    .AddInteractiveServerComponents();
builder.Services.AddRazorPages();

var app = builder.Build();

if (app.Environment.IsDevelopment())
{
    app.UseSwagger();
    app.UseSwaggerUI();
}

app.UseStaticFiles();
app.UseRouting();

app.UseAntiforgery();
app.MapRazorPages();

app.MapRazorComponents<App>()
    .AddInteractiveServerRenderMode();
app.MapControllers();
app.Run();