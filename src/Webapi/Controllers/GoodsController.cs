using Azure.Storage.Blobs;
using Microsoft.AspNetCore.Mvc;
using Microsoft.EntityFrameworkCore;
using mmp.Data;
using ZiggyCreatures.Caching.Fusion;
using SixLabors.ImageSharp;
using SixLabors.ImageSharp.Processing;
using SixLabors.ImageSharp.Formats.Jpeg;


namespace Webapi.Controllers
{
    [Route("api/goods")]
    [ApiController]
    public class GoodsController : ControllerBase
    {
        private readonly ApplicationDbContext db;
        private readonly IHostEnvironment host;
        private readonly IFusionCache cache;

        private readonly BlobContainerClient? blobContainer;

        private bool UseAzureBlobs => !host.IsDevelopment();

        private static readonly string[] cacheTags = ["goodImages"];

        private void ClearCache() => cache.RemoveByTag(cacheTags);

        public GoodsController(ApplicationDbContext _db, IHostEnvironment hostEnvironment, BlobServiceClient _blobServiceClient, IFusionCache _cache)
        {
            db = _db;
            host = hostEnvironment;
            cache = _cache;

            if (UseAzureBlobs && _blobServiceClient != null)
                blobContainer = _blobServiceClient.GetBlobContainerClient("goodimages");
        }

        [HttpGet]
        public async Task<ActionResult<IEnumerable<Good>>> GetGoods([FromQuery] long shopId)
        {
            return await db.Goods
                .AsNoTracking()
                .Where(_ => _.OwnerShop.ID == shopId && !_.IsDeleted)
                .OrderBy(_ => _.OrderInShop)
                .ThenBy(_ => _.Caption)
                .ToListAsync();
        }

        [HttpGet("{id}")]
        public async Task<ActionResult<Good>> GetGood(long id, [FromQuery] bool showDeleted = false)
        {
            var good = await db.Goods
                .Include(_ => _.OwnerShop)
                .AsNoTracking()
                .FirstOrDefaultAsync(_ => _.ID == id && (!_.IsDeleted || showDeleted));
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
            dbGood.OrderInShop = good.OrderInShop;

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
                Url = good.Url,
                OrderInShop = good.OrderInShop
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

            for (var i = 0; i < 3; i++)
            {
                var r0 = await DeleteGoodImage(id, i);
                if (r0 is StatusCodeResult r)
                    Console.WriteLine($"Deleting good {id}. image {i} deleting result {r.StatusCode}");
                else
                    Console.WriteLine($"Deleting good {id}. image {i}. response type is {r0.GetType().ToString()}");
            }

            good.IsDeleted = true;

            // also delete uncompleted baskets
            db.OrderLines.Where(_ => _.Good == good && _.Order == null).ExecuteDelete();

            await db.SaveChangesAsync();

            return NoContent();
        }

        #region images

        private static string BlobName(long goodId, int imgNum) => $"goodImage-{goodId}-{imgNum}";
        private const string DevBlobStorageFolder = "c:\\tmp\\";

        [HttpGet("images/{goodId:long}/{num:int}")]
        public async Task<IActionResult> GetGoodImage(long goodId, int num)
        {
            var res = await cache.GetOrSetAsync(
                $"goodImage:{goodId}:{num}",
                async (ctx) =>
                {
                    if (!UseAzureBlobs)
                    {
                        var fs = DevBlobStorageFolder + BlobName(goodId, num);
                        if (!System.IO.File.Exists(fs)) return null;
                        using var fstream = new FileStream(fs, FileMode.Open);
                        using var mstream = new MemoryStream();
                        fstream.CopyTo(mstream);
                        return File(mstream.ToArray(), "image/jpeg");
                    }
                    else
                    {
                        if (blobContainer == null)
                            throw new Exception("BlobContainer=null. Check blob storage credentials");
                        var handler = blobContainer.GetBlobClient(BlobName(goodId, num));
                        if (!handler.Exists()) return null;
                        using var memoryStream = new MemoryStream();
                        await handler.DownloadToAsync(memoryStream);
                        memoryStream.Position = 0;
                        return File(memoryStream.ToArray(), "image/jpeg");
                    }
                },
                tags: cacheTags
            );

            if (res == null) return NoContent();
            return res;
        }

        [HttpGet("thumbs/{goodId:long}/{num:int}")]
        public async Task<IActionResult> GetGoodThumb(long goodId, int num)
        {
            var res = await cache.GetOrSetAsync(
                $"goodThumb:{goodId}:{num}",
                async (ctx) =>
                {
                    var res = await GetGoodImage(goodId, num);
                    if (res is NoContentResult) return null;
                    if (res is not FileContentResult imgCnt) return null;

                    var ms = new MemoryStream(imgCnt.FileContents);

                    var img = Image.Load(ms);
                    img.Mutate(x => x.Resize(60, 60));
                    ms = new MemoryStream();

                    img.Save(ms, new JpegEncoder());
                    ms.Position = 0;
                    return File(ms.ToArray(), "image/jpeg");

                },
                tags: cacheTags
            );
            if (res == null) return NoContent();
            return res;
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
                if (blobContainer == null)
                    throw new Exception("BlobContainer=null. Check blob storage credentials");
                await blobContainer.DeleteBlobIfExistsAsync(BlobName(goodId, num));
                var handler = blobContainer.GetBlobClient(BlobName(goodId, num));
                using var memoryStream = new MemoryStream();
                image.CopyTo(memoryStream);
                memoryStream.Position = 0;
                await handler.UploadAsync(memoryStream);
            }

            ClearCache();
            return NoContent();
        }

        [HttpDelete("images/{goodId:long}/{num:int}")]
        public async Task<IActionResult> DeleteGoodImage(long goodId, int num)
        {
            var cu = db.CurrentUser();
            if (cu == null) return Unauthorized("Unable to get authenticated user");

            var good = await db.Goods.
                Include(_ => _.OwnerShop)
                .FirstOrDefaultAsync(_ => _.ID == goodId && !_.IsDeleted);
            if (good == null) return NotFound("Good isn't found");
            if (good.OwnerShop.CreatedByID != cu.Id) return Unauthorized("Good isn't your. Only owner can delete");

            if (!UseAzureBlobs)
            {
                var fs = DevBlobStorageFolder + BlobName(goodId, num);
                if (System.IO.File.Exists(fs)) System.IO.File.Delete(fs);
            }
            else
            {
                if (blobContainer == null)
                    throw new Exception("BlobContainer=null. Check blob storage credentials");
                await blobContainer.DeleteBlobIfExistsAsync(BlobName(goodId, num));
            }

            ClearCache();
            return NoContent();
        }

        #endregion

    }
}
