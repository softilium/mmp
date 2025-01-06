using Microsoft.AspNetCore.Authorization;
using Microsoft.AspNetCore.Mvc;
using Microsoft.EntityFrameworkCore;
using mmp.DbCtx;
using mmp.Models;
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

        [HttpGet("outbox")]
        [Authorize]
        public async Task<ActionResult<IEnumerable<Order>>> GetOrders([FromQuery] int showAll)
        {
            var cu = db.CurrentUser();
            if (cu == null) return Unauthorized();

            return await db.Orders
                .Where(_ => !_.IsDeleted && _.CreatedBy.Id == cu.Id)
                .Where(_ => showAll == 1 || (_.Status != OrderStatuses.Done && _.Status != OrderStatuses.Canceled))
                .Include(_ => _.Shop)
                .OrderByDescending(_ => _.CreatedOn)
                .ToListAsync();
        }

        [HttpGet("outbox/{id}")]
        public async Task<ActionResult<Order>> GetOrder(long id)
        {
            var cu = db.CurrentUser();
            if (cu == null) return Unauthorized();

            var order = await db.Orders
                .Include(_=>_.Shop)
                .Include(_=>_.Lines)
                    .ThenInclude(_=>_.Good)
                .FirstOrDefaultAsync(_ => _.ID == id && !_.IsDeleted && _.CreatedBy.Id == cu.Id);
            if (order == null) return NotFound();

            return order;
        }

        [HttpPut("outbox/{id}")]
        [Authorize]
        public async Task<IActionResult> PutOrder(long id, Order order)
        {
            var dborder = await GetOrder(id);
            if (dborder.Value == null) return NotFound();

            dborder.Value.Status = order.Status;
            await db.SaveChangesAsync();
            return NoContent();
        }

        [HttpPost("outbox/{shopid:long}")]
        [Authorize]
        public async Task<ActionResult<Order>> PostOrder([FromRoute] long shopId)
        {
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
                Sum = lines.Sum(_ => _.Sum)
            };
            db.Orders.Add(dborder);
            foreach (var line in lines) line.Order = dborder;

            await db.SaveChangesAsync();
            return CreatedAtAction("GetOrder", new { id = dborder.ID }, dborder);
        }

    }
}
