using Creastina.CraftingTool.Models;

namespace Creastina.CraftingTool.Repository;

public interface ITodoRepository
{
    public Task<List<Todo>> GetAllTodos();
    public Task<Todo?> GetTodoById(int id);
    public Task<Todo> CreateTodo(Todo todo);
    public Task UpdateTodo(Todo todo);
    public Task DeleteTodo(int id);
}