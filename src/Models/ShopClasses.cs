using System.ComponentModel.DataAnnotations;
using Microsoft.AspNetCore.Identity;
namespace mmp.Models
{
    public class User : IdentityUser<long> { }

    public abstract class BaseObject
    {
        [Key] public long ID { get; set; }

        [Required] public long CreatedByID { get; set; }
        [Required] public DateTime CreatedOn { get; set; }

        [Required] public long ModifiedByID { get; set; }
        [Required] public DateTime ModifiedOn { get; set; }

        public bool IsDeleted { get; set; } = false;
        public long? DeletedByID { get; set; }
        public DateTime? DeletedOn { get; set; }

        public string? Comment { get; set; }

        public void BeforeSave() { }
    }

    public class Shop : BaseObject
    {
        [Required] public long OwnerUserID { get; set; }
        [Required] public string Name { get; set; } = "shop1";
        [Required] public string Caption { get; set; } = "Shop 1";

    }

    public class Good : BaseObject
    {
        [Required] public long OwnerShopID { get; set; }
        [Required] public string Caption { get; set; } = "";
        public string? Description { get; set; }
    }

    public enum OrderStatuses : int
    {
        New = 100,
        ReadyToDeliver = 200,
        Delivered = 300,
        Canceled = 400,
    }

    public class Order : BaseObject
    {
        [Required] public long ShopID { get; set; }
        [Required] public long BuyerID { get; set; }
        [Required] public OrderStatuses Status { get; set; }
    }

    public class OrderLine : BaseObject
    {
        [Required] public long WhoID { get; set; }
        [Required] public long ShopID { get; set; }     
        public long? OrderID { get; set; } //when empty it is basket
        [Required] public long GoodID { get; set; }
        [Required] public long Qty { get; set; }
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
