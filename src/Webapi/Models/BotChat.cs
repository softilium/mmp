using System.ComponentModel.DataAnnotations;

namespace mmp.Models
{

    public class BotChat
    {
        [Key]
        [Required]
        public long ChatId { get; set; }

        [Required]
        [MaxLength(50)]
        public string UserName { get; set; }
    }

}
