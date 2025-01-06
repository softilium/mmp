using System.ComponentModel;
using System.ComponentModel.DataAnnotations;
using Microsoft.AspNetCore.Identity;
using Microsoft.EntityFrameworkCore;

namespace mmp.Models
{
    public class User : IdentityUser<long> { }

    public abstract class BaseObject
    {
        [Key] public long ID { get; set; }

        [Required][DeleteBehavior(DeleteBehavior.Restrict)] public User CreatedBy { get; set; }
        [Required] public DateTime CreatedOn { get; set; }

        [DeleteBehavior(DeleteBehavior.Restrict)] public User? ModifiedBy { get; set; }
        public DateTime ModifiedOn { get; set; }

        public bool IsDeleted { get; set; } = false;
        public User? DeletedBy { get; set; }
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

//todo tags
//todo categories
//todo Prices columns
//todo integration-telegram-payments
//todo comments
//todo search (elastic)
//todo telegram for notify
//todo blacklist
//todo voting for goods/shops (rating, ...)
//todo recommendations
//todo location for shops/delivery
//todo want.to.buy.later (notify me when price go low)
//todo notify all past customers
//todo docker
//todo web blazor client
//todo maui ios
//todo maui android
//todo telegram webapp
//todo azure sql
//todo azure funcs
//todo azure app services
//todo azure static web apps
//todo xunit, nunit
//todo typescript
//todo redis
//todo llm
//todo mini-profiler
//todo opened-closed shops (public, direct-link, approved-only)
//todo dataversion and checking for it before update

//todo good fields: url, images
//todo images storage (guid+file, separated collection for each Good)
//todo use special proxy for IdentityUser
//todo goods page
//todo profile page
//todo incoming orders page
//todo email/telegram notify
//todo stop to use email globally inside vueclient/api
