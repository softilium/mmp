using System.ComponentModel.DataAnnotations;

namespace mmp.Data
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

    //todo clean old records. When user verify it, we put chatId value into user record

}
