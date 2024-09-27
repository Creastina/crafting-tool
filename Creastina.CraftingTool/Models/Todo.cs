using System.Text.Json.Serialization;

namespace Creastina.CraftingTool.Models;

public record Todo
{
    [JsonIgnore] public int Id { get; set; }

    [JsonPropertyName("id")] public int ReadId => Id;

    public string Title { get; set; }
    public string? Status { get; set; }
    public string? Kind { get; set; }
    public string? Material { get; set; }
    public bool IsDone { get; set; }
    public bool IsPartsMissing { get; set; }
    public string? Notes { get; set; }
}