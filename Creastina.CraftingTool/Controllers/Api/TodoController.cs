using System.ComponentModel.DataAnnotations;
using Creastina.CraftingTool.Authentication;
using Creastina.CraftingTool.Models;
using Creastina.CraftingTool.Repository;
using Microsoft.AspNetCore.Mvc;

namespace Creastina.CraftingTool.Controllers.Api;

[ApiController]
[ApiKey]
[Route("api/todo")]
public class TodoController(ITodoRepository todoRepository) : ControllerBase
{
    [HttpGet]
    [Produces("application/json")]
    [ProducesResponseType(StatusCodes.Status200OK)]
    public async Task<List<Todo>> GetAll()
    {
        return await todoRepository.GetAllTodos();
    }

    [HttpGet("{id:int}")]
    [Produces("application/json")]
    [ProducesResponseType<Todo>(StatusCodes.Status200OK)]
    [ProducesResponseType(StatusCodes.Status404NotFound)]
    public async Task<IActionResult> GetById(int id)
    {
        var todo = await todoRepository.GetTodoById(id);
        if (todo != null) return Ok(todo);

        return NotFound();
    }

    [HttpPost]
    [Produces("application/json")]
    [ProducesResponseType<Todo>(StatusCodes.Status201Created)]
    [ProducesResponseType(StatusCodes.Status500InternalServerError)]
    [Consumes("application/json")]
    public async Task<IActionResult> Create([FromBody] [Required] Todo todo)
    {
        var created = await todoRepository.CreateTodo(todo);

        return Created(Url.ActionLink("GetById", "Todo", new { id = created.Id }), created);
    }

    [HttpPut("{id:int}")]
    [Produces("application/json")]
    [ProducesResponseType<Todo>(StatusCodes.Status204NoContent)]
    [ProducesResponseType(StatusCodes.Status404NotFound)]
    [ProducesResponseType(StatusCodes.Status500InternalServerError)]
    [Consumes("application/json")]
    public async Task<IActionResult> Update(int id, [FromBody] [Required] Todo todo)
    {
        try
        {
            todo.Id = id;
            await todoRepository.UpdateTodo(todo);

            return NoContent();
        }
        catch (EntryNotFoundException)
        {
            return NotFound();
        }
    }

    [HttpDelete("{id:int}")]
    [Produces("application/json")]
    [ProducesResponseType<Todo>(StatusCodes.Status204NoContent)]
    [ProducesResponseType(StatusCodes.Status404NotFound)]
    [ProducesResponseType(StatusCodes.Status500InternalServerError)]
    [Consumes("application/json")]
    public async Task<IActionResult> Delete(int id)
    {
        try
        {
            await todoRepository.DeleteTodo(id);

            return NoContent();
        }
        catch (EntryNotFoundException)
        {
            return NotFound();
        }
    }
}