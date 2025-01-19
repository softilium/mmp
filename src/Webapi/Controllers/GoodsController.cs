using Azure.Storage.Blobs;
using Microsoft.AspNetCore.Mvc;
using Microsoft.EntityFrameworkCore;
using mmp.Data;

namespace Webapi.Controllers
{
    [Route("api/goods")]
    [ApiController]
    public class GoodsController : ControllerBase
    {
        private readonly ApplicationDbContext db;
        private IHostEnvironment host;

        private readonly BlobContainerClient blobContainer;

        private bool UseAzureBlobs => !host.IsDevelopment();

        public GoodsController(ApplicationDbContext _db, IHostEnvironment hostEnvironment, BlobServiceClient _blobServiceClient)
        {
            db = _db;
            host = hostEnvironment;
            if (UseAzureBlobs && _blobServiceClient != null)
                blobContainer = _blobServiceClient.GetBlobContainerClient("goodimages");
        }

        [HttpGet]
        public async Task<ActionResult<IEnumerable<Good>>> GetGoods([FromQuery] long shopId)
        {
            return await db.Goods.AsNoTracking().Where(_ => _.OwnerShop.ID == shopId && !_.IsDeleted).ToListAsync();
        }

        [HttpGet("{id}")]
        public async Task<ActionResult<Good>> GetGood(long id)
        {
            var good = await db.Goods.Include(_ => _.OwnerShop).AsNoTracking().FirstOrDefaultAsync(_ => _.ID == id);
            if (good == null) return NotFound();
            return good;
        }

        [HttpPut("{id}")]
        public async Task<IActionResult> PutGood(long id, Good good)
        {
            if (id != good.ID) return BadRequest();

            var cu = db.CurrentUser();
            if (cu == null) return Unauthorized();

            var dbGood = await db.Goods.FirstOrDefaultAsync(_ => _.ID == id && !_.IsDeleted);
            if (dbGood == null) return NotFound();
            if (dbGood.CreatedByID != cu.Id) return Unauthorized();

            dbGood.Caption = good.Caption;
            dbGood.Description = good.Description;
            dbGood.Price = good.Price;
            dbGood.Article = good.Article;
            dbGood.Url = good.Url;

            await db.SaveChangesAsync();

            return NoContent();
        }

        [HttpPost]
        public async Task<ActionResult<Good>> PostGood(Good good)
        {
            var cu = db.CurrentUser();
            if (cu == null) return Unauthorized();

            var shop = await db.Shops.FirstOrDefaultAsync(_ => _.ID == good.OwnerShop.ID && !_.IsDeleted);
            if (shop == null) return NotFound();

            if (shop.CreatedByID != cu.Id) return Unauthorized();

            var dbGood = new Good
            {
                OwnerShop = shop,
                Caption = good.Caption,
                Description = good.Description,
                Price = good.Price,
                Article = good.Article,
                Url = good.Url
            };

            db.Goods.Add(dbGood);
            await db.SaveChangesAsync();

            return CreatedAtAction("GetGood", new { id = dbGood.ID }, dbGood);
        }

        [HttpDelete("{id}")]
        public async Task<IActionResult> DeleteGood(long id)
        {
            var cu = db.CurrentUser();
            if (cu == null) return Unauthorized();

            var good = await db.Goods.FirstOrDefaultAsync(_ => _.ID == id && !_.IsDeleted);
            if (good == null) return NotFound();
            if (good.CreatedByID != cu.Id) return Unauthorized();

            good.IsDeleted = true;

            // also delete uncompleted baskets
            db.OrderLines.Where(_ => _.Good == good && _.Order == null).ExecuteDelete();

            await DeleteGoodImage(id, 1);
            await DeleteGoodImage(id, 2);
            await DeleteGoodImage(id, 3);

            await db.SaveChangesAsync();

            return NoContent();
        }

        #region images

        private static string BlobName(long goodId, int imgNum) => $"goodImage-{goodId}-{imgNum}";
        private const string DevBlobStorageFolder = "c:\\tmp\\";

        [HttpGet("images/{goodId:long}/{num:int}")]
        public async Task<IActionResult> GetGoodImage(long goodId, int num)
        {
            if (!UseAzureBlobs)
            {
                var fs = DevBlobStorageFolder + BlobName(goodId, num);
                if (!System.IO.File.Exists(fs)) return NoContent();
                return PhysicalFile(fs, "image/jpeg");
            }
            else
            {
                var handler = blobContainer.GetBlobClient(BlobName(goodId, num));
                if (!handler.Exists()) return NoContent();
                using var memoryStream = new MemoryStream();
                await handler.DownloadToAsync(memoryStream);
                memoryStream.Position = 0;
                return File(memoryStream.ToArray(), "image/jpeg");
            }
        }

        [HttpPost("images/{goodId:long}/{num:int}")]
        public async Task<IActionResult> PostGoodImage(long goodId, int num)
        {
            var image = Request.Form.Files.Count > 0 ? Request.Form.Files[0] : null;
            if (image == null || image.Length == 0) return BadRequest();

            var cu = db.CurrentUser();
            if (cu == null) return Unauthorized();

            var good = await db.Goods.
                Include(_ => _.OwnerShop)
                .FirstOrDefaultAsync(_ => _.ID == goodId && !_.IsDeleted);
            if (good == null) return NotFound();
            if (good.OwnerShop.CreatedByID != cu.Id) return Unauthorized();

            if (!UseAzureBlobs)
            {
                using var stream = System.IO.File.Create(DevBlobStorageFolder + BlobName(goodId, num));
                image.CopyTo(stream);
            }
            else
            {
                await blobContainer.DeleteBlobIfExistsAsync(BlobName(goodId, num));
                var handler = blobContainer.GetBlobClient(BlobName(goodId, num));
                using var memoryStream = new MemoryStream();
                image.CopyTo(memoryStream);
                memoryStream.Position = 0;
                await handler.UploadAsync(memoryStream);
            }
            return NoContent();
        }

        [HttpDelete("images/{goodId:long}/{num:int}")]
        public async Task<IActionResult> DeleteGoodImage(long goodId, int num)
        {
            var cu = db.CurrentUser();
            if (cu == null) return Unauthorized();

            var good = await db.Goods.
                Include(_ => _.OwnerShop)
                .FirstOrDefaultAsync(_ => _.ID == goodId && !_.IsDeleted);
            if (good == null) return NotFound();
            if (good.OwnerShop.CreatedByID != cu.Id) return Unauthorized();

            if (!UseAzureBlobs)
            {
                var fs = DevBlobStorageFolder + BlobName(goodId, num);
                if (System.IO.File.Exists(fs)) System.IO.File.Delete(fs);
            }
            else
            {
                await blobContainer.DeleteBlobIfExistsAsync(BlobName(goodId, num));
            }
            return NoContent();
        }

        #endregion

    }
}
