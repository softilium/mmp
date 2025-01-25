using Microsoft.EntityFrameworkCore;
using System.ComponentModel.DataAnnotations;

namespace mmp.Data
{
    public class Good : BaseObject
    {
        [Required]
        [DeleteBehavior(DeleteBehavior.Restrict)] 
        public Shop OwnerShop { get; set; }
        
        [MaxLength(100)]
        [Required] 
        public string Caption { get; set; } = "";
        
        [MaxLength(50)] 
        public string? Article { get; set; }
        
        [MaxLength(900)] 
        public string? Url { get; set; }
        
        public string? Description { get; set; }
        
        [Precision(15, 2)] 
        public decimal Price { get; set; }
        
        public int OrderInShop { get; set; } = 100;
    }
}
