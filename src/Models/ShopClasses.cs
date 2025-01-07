using System.ComponentModel;
using System.ComponentModel.DataAnnotations;
using System.ComponentModel.DataAnnotations.Schema;
using Microsoft.AspNetCore.Identity;
using Microsoft.EntityFrameworkCore;

namespace mmp.Models
{
    // User's projection for frontend
    public class UserInfo
    {
        public string UserName { get; set; }
        public bool ShopManage { get; set; }
        public bool Admin { get; set; }

        public UserInfo(User src)
        {
            ArgumentNullException.ThrowIfNull(src);
            UserName = src.UserName ?? "";
            ShopManage = src.ShopManage;
            Admin = src.Admin;
        }
    }

    public class User : IdentityUser<long>
    {
        public bool ShopManage { get; set; }
        public bool Admin { get; set; }
    }

    [Index(nameof(CreatedByID))]
    public abstract class BaseObject
    {
        [Key] public long ID { get; set; }

        [Required]
        public long CreatedByID { get; set; }
        [NotMapped] public UserInfo? CreatedByInfo { get; set; }
        [Required] public DateTime CreatedOn { get; set; }

        public long? ModifiedByID { get; set; }
        [NotMapped] public UserInfo? User { get; set; }
        public DateTime ModifiedOn { get; set; }

        public bool IsDeleted { get; set; } = false;
        public long? DeletedByID { get; set; }
        [NotMapped] public UserInfo? DeletedByInfo { get; set; }
        public DateTime? DeletedOn { get; set; }

        public void BeforeSave() { }
    }

    public class Shop : BaseObject
    {
        [Required] public string Caption { get; set; } = "Shop 1";
    }

    public class Good : BaseObject
    {
        [Required][DeleteBehavior(DeleteBehavior.Restrict)] public Shop OwnerShop { get; set; }
        [Required] public string Caption { get; set; } = "";
        public string? Description { get; set; }
        [Precision(15, 2)] public decimal Price { get; set; }
    }

    public enum OrderStatuses : int
    {
        [Description("Новый")]
        New = 100,

        [Description("В работе")]
        InProcess = 200,

        [Description("Готов к доставке")]
        ReadyToDeliver = 300,

        [Description("Доставляется")]
        Delivering = 400,

        [Description("Доставлен")]
        Done = 500,

        [Description("Отменен")]
        Canceled = 600,
    }

    public class Order : BaseObject
    {
        [Required][DeleteBehavior(DeleteBehavior.Restrict)] public Shop Shop { get; set; }
        [Required] public OrderStatuses Status { get; set; }
        [Required][Precision(15, 2)] public decimal Qty { get; set; }
        [Required][Precision(15, 2)] public decimal Sum { get; set; }
        public ICollection<OrderLine> Lines { get; set; } = [];
    }

    public class OrderLine : BaseObject
    {
        [Required][DeleteBehavior(DeleteBehavior.Restrict)] public Shop Shop { get; set; }
        public Order? Order { get; set; } //when empty it is basket
        [Required][DeleteBehavior(DeleteBehavior.Restrict)] public Good Good { get; set; }
        [Required][Precision(15, 2)] public decimal Qty { get; set; }
        [Required][Precision(15, 2)] public decimal Sum { get; set; }
    }
}

