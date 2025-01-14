using System.ComponentModel.DataAnnotations;

namespace mmp.Models
{
    public class Shop : BaseObject
    {
        [MaxLength(100)]
        [Required] 
        public string Caption { get; set; } = "Shop 1";
    }
}
