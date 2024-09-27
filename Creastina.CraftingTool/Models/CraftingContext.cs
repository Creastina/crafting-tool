using Microsoft.EntityFrameworkCore;

namespace Creastina.CraftingTool.Models;

public class CraftingContext(DbContextOptions<CraftingContext> options) : DbContext(options)
{
    public DbSet<Todo> Todos { get; set; }
}