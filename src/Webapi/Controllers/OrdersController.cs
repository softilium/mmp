using Microsoft.AspNetCore.Mvc;
using Microsoft.EntityFrameworkCore;
using mmp.Data;
using System.ComponentModel;

namespace Webapi.Controllers
{

    public static class EnumExtensionMethods
    {
        public static string GetEnumDescription(this OrderStatuses enumValue)
        {
            var fieldInfo = enumValue.GetType().GetField(enumValue.ToString());
            if (fieldInfo == null) return "";

            var descriptionAttributes = (DescriptionAttribute[])fieldInfo.GetCustomAttributes(typeof(DescriptionAttribute), false);

            return descriptionAttributes.Length > 0 ? descriptionAttributes[0].Description : enumValue.ToString();
        }
    }

    [Route("api/orders")]
    [ApiController]
    public class OrdersController : ControllerBase
    {
        private readonly ApplicationDbContext db;

        public OrdersController(ApplicationDbContext _db)
        {
            db = _db;
        }

        [HttpGet("statuses")]
        public Dictionary<int, string> GetStatuses()
        {
            Dictionary<int, string> res = [];
            foreach (var st in Enum.GetValues<OrderStatuses>())
                res.Add((int)st, st.GetEnumDescription());
            return res;
        }

        [HttpGet("{id}")]
        public async Task<ActionResult<Order>> GetOrder(long id)
        {
            var cu = db.CurrentUser();
            if (cu == null) return Unauthorized();

            var order = await db.Orders
                .Include(_ => _.Lines)
                    .ThenInclude(_ => _.Good)
                .FirstOrDefaultAsync(_ => _.ID == id && !_.IsDeleted);
            if (order == null) return NotFound();

            if (order.CreatedByID != cu.Id && order.Shop.CreatedByID != cu.Id) return Unauthorized();

            return order;
        }


        [HttpGet("outbox")]
        public async Task<ActionResult<IEnumerable<Order>>> GetOrders([FromQuery] int showAll)
        {
            var cu = db.CurrentUser();
            if (cu == null) return Unauthorized();

            var r = await db.Orders
                .Where(_ => !_.IsDeleted && _.CreatedByID == cu.Id)
                .Where(_ => showAll == 1 || (_.Status != OrderStatuses.Done && _.Status != OrderStatuses.Canceled))
                .OrderByDescending(_ => _.CreatedOn)
                .ToListAsync();
            foreach (var item in r) UserCache.LoadCreatedBy(item.Shop, db);
            return r;
        }

        [HttpPut("outbox/{id:long}")]
        public async Task<IActionResult> PutCustomerOrder(long id, Order order)
        {
            var cu = db.CurrentUser();
            if (cu == null) return Unauthorized();

            var dborder = db.Orders.FirstOrDefault(_ => _.ID == id);
            if (dborder == null) return NotFound();

            if (order.CreatedByID != cu.Id) return Unauthorized();

            dborder.CustomerComment = order.CustomerComment;

            await db.SaveChangesAsync();
            return NoContent();
        }

        [HttpPost("outbox/{shopid:long}")]
        public async Task<ActionResult<Order>> PostCustomerOrder([FromRoute] long shopId)
        {
            using StreamReader reader = new(Request.Body, leaveOpen: true);
            string customerComment = await reader.ReadToEndAsync();

            var cu = db.CurrentUser();
            if (cu == null) return Unauthorized();

            var shop = db.Shops.FirstOrDefault(_ => !_.IsDeleted && _.ID == shopId);
            if (shop == null) return NotFound();

            var lines = db.OrderLines.Where(_ => _.Order == null && _.Shop.ID == shopId).ToList();
            if (lines.Count == 0) return NotFound();

            var dborder = new Order
            {
                Shop = shop,
                Status = OrderStatuses.New,
                Qty = lines.Sum(_ => _.Qty),
                Sum = lines.Sum(_ => _.Sum),
                CustomerComment = customerComment
            };
            db.Orders.Add(dborder);
            foreach (var line in lines) line.Order = dborder;

            await db.SaveChangesAsync();
            return CreatedAtAction("GetOrder", new { id = dborder.ID }, dborder);
        }

        [HttpGet("inbox")]
        public async Task<ActionResult<IEnumerable<Order>>> GetSenderOrders([FromQuery] int showAll)
        {
            var cu = db.CurrentUser();
            if (cu == null) return Unauthorized();

            var myshops = db.Shops.Where(_ => !_.IsDeleted && _.CreatedByID == cu.Id).Select(_ => _.ID).ToList();

            var r = await db.Orders
                .Where(_ => !_.IsDeleted && myshops.Contains(_.Shop.ID))
                .Where(_ => showAll == 1 || (_.Status != OrderStatuses.Done && _.Status != OrderStatuses.Canceled))
                .OrderByDescending(_ => _.CreatedOn)
                .ToListAsync();
            foreach (var item in r)
            {
                UserCache.LoadCreatedBy(item, db);
                UserCache.LoadCreatedBy(item.Shop, db);
            }
            return r;
        }

        [HttpPut("inbox/{id:long}")]
        public async Task<ActionResult<IEnumerable<Order>>> PutSenderOrder(long id, Order order)
        {
            var cu = db.CurrentUser();
            if (cu == null) return Unauthorized();

            var dborder = db.Orders.FirstOrDefault(_ => _.ID == id);
            if (dborder == null) return NotFound();

            if (cu.Id != order.Shop.CreatedByID) return Unauthorized();

            dborder.Status = order.Status;
            dborder.SenderComment = order.SenderComment;
            dborder.ExpectedDeliveryDate = order.ExpectedDeliveryDate;

            await db.SaveChangesAsync();

            return NoContent();
        }
    }
}
