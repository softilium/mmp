using Microsoft.EntityFrameworkCore;
using System.ComponentModel.DataAnnotations.Schema;
using System.ComponentModel.DataAnnotations;

namespace mmp.Models
{
    [Index(nameof(CreatedByID))]
    public abstract class BaseObject
    {
        [Key] 
        public long ID { get; set; }

        [Required]
        public long CreatedByID { get; set; }
        
        [NotMapped] 
        public UserInfo? CreatedByInfo { get; set; }
        
        [Required] 
        public DateTime CreatedOn { get; set; }

        public long? ModifiedByID { get; set; }
        
        [NotMapped] 
        public UserInfo? User { get; set; }
        public DateTime ModifiedOn { get; set; }

        public bool IsDeleted { get; set; } = false;
        public long? DeletedByID { get; set; }
        
        [NotMapped] 
        public UserInfo? DeletedByInfo { get; set; }
        public DateTime? DeletedOn { get; set; }

        public void BeforeSave() { }
        public void AfterSave() { }
    }
}
