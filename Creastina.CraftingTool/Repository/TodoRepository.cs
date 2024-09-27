using Creastina.CraftingTool.Models;
using Microsoft.EntityFrameworkCore;

namespace Creastina.CraftingTool.Repository;

public class TodoRepository(CraftingContext context) : ITodoRepository
{
    public Task<List<Todo>> GetAllTodos()
    {
        return context.Todos.ToListAsync();
    }

    public Task<Todo?> GetTodoById(int id)
    {
        return context.FindAsync<Todo>(id).AsTask();
    }

    public async Task<Todo> CreateTodo(Todo todo)
    {
        var newTodo = await context.AddAsync(todo);
        await context.SaveChangesAsync();

        return newTodo.Entity;
    }

    public async Task UpdateTodo(Todo todo)
    {
        if (await context.Todos.AnyAsync(t => t.Id == todo.Id))
            await context.Todos.Where(t => t.Id == todo.Id)
                .ExecuteUpdateAsync(calls =>
                    calls
                        .SetProperty(t => t.Kind, todo.Kind)
                        .SetProperty(t => t.Material, todo.Material)
                        .SetProperty(t => t.Title, todo.Title)
                        .SetProperty(t => t.Notes, todo.Notes)
                        .SetProperty(t => t.Status, todo.Status)
                        .SetProperty(t => t.IsDone, todo.IsDone)
                        .SetProperty(t => t.IsPartsMissing, todo.IsPartsMissing)
                );
        else
            throw new EntryNotFoundException();
    }

    public async Task DeleteTodo(int id)
    {
        if (await context.Todos.AnyAsync(t => t.Id == id))
            await context.Todos.Where(todo => todo.Id == id).ExecuteDeleteAsync();
        else
            throw new EntryNotFoundException();
    }
}