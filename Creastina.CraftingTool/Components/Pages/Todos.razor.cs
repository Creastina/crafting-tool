using Blazored.LocalStorage;
using Creastina.CraftingTool.Authentication;
using Creastina.CraftingTool.Models;
using Creastina.CraftingTool.Repository;
using Microsoft.AspNetCore.Components;
using Microsoft.JSInterop;

namespace Creastina.CraftingTool.Components.Pages;

public partial class Todos : ComponentBase
{
    [Inject] ITodoRepository TodoRepository { get; set; }
    [Inject] IJSRuntime JsRuntime { get; set; }
    [Inject] ILocalStorageService LocalStorageService { get; set; }
    [Inject] NavigationManager NavigationManager { get; set; }
    [Inject] CookieAuthenticationChecker CookieAuthenticationChecker { get; set; }
    [Inject] IHttpContextAccessor HttpContextAccessor { get; set; }

    private List<Todo> _todos = [];

    private Todo _selectedTodo = new();

    protected override async Task OnParametersSetAsync()
    {
        CookieAuthenticationChecker.CheckCookie();

        _todos = await TodoRepository.GetAllTodos();
        _selectedTodo = _todos.FirstOrDefault() ?? new Todo { IsNew = true };

        await base.OnParametersSetAsync();
    }

    private async Task SaveTodo()
    {
        CookieAuthenticationChecker.CheckCookie();

        if (_selectedTodo.IsNew)
        {
            _selectedTodo = await TodoRepository.CreateTodo(_selectedTodo);
            _selectedTodo.IsNew = false;
            _todos = await TodoRepository.GetAllTodos();
        }
        else
        {
            await TodoRepository.UpdateTodo(_selectedTodo);
            _todos = await TodoRepository.GetAllTodos();
            _selectedTodo = _todos.SingleOrDefault(todo => todo.Id == _selectedTodo.Id) ?? new Todo { IsNew = true };
        }
    }

    private async Task DeleteTodo()
    {
        CookieAuthenticationChecker.CheckCookie();

        await TodoRepository.DeleteTodo(_selectedTodo.Id);
        _todos = await TodoRepository.GetAllTodos();
        _selectedTodo = _todos.FirstOrDefault() ?? new Todo { IsNew = true };
        await JsRuntime.InvokeVoidAsync("confirmDelete.close");
    }
}