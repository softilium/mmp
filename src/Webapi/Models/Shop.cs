using System.ComponentModel.DataAnnotations;

namespace mmp.Data
{
    public class Shop : BaseObject
    {
        [MaxLength(100)]
        [Required] 
        public string Caption { get; set; } = "Shop 1";

        [MaxLength(300)]
        public string Description { get; set; } = "";

        [MaxLength(300)]
        public string DeliveryConditions { get; set; } = "";

    }
}
