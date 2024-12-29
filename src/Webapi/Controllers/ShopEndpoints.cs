using Microsoft.EntityFrameworkCore;
using Microsoft.AspNetCore.Http.HttpResults;
using mmp.DbCtx;
using mmp.Models;

namespace mmp.Api
{

    public static class ShopEndpoints
    {
        public static void MapShopEndpoints(this IEndpointRouteBuilder routes)
        {
            var group = routes.MapGroup("/api/shops");

            group.MapGet("/", async (ApplicationDbContext db) =>
            {
                return await db.Shops.ToListAsync();
            })
            .WithName("GetAllShops")
            .WithOpenApi();

            group.MapGet("/{id:long}", async Task<Results<Ok<Shop>, NotFound>> (long id, ApplicationDbContext db) =>
            {
                return await db.Shops.Where(_ => !_.IsDeleted).AsNoTracking()
                    .FirstOrDefaultAsync(model => model.ID == id)
                    is Shop model
                        ? TypedResults.Ok(model)
                        : TypedResults.NotFound();
            })
            .WithName("GetShopById")
            .WithOpenApi();

            group.MapPut("/{id:long}", async Task<Results<Ok, NotFound>> (long id, Shop shop, ApplicationDbContext db) =>
            {
                var dbobj = db.Shops.First(_ => _.ID == id);
                if (dbobj == null) return TypedResults.NotFound();

                dbobj.OwnerUserID = shop.OwnerUserID;
                dbobj.Name = shop.Name;
                dbobj.Caption = shop.Caption;
                dbobj.Comment = shop.Comment;

                await db.SaveChangesAsync();

                return TypedResults.Ok();
            })
            .WithName("UpdateShop")
            .WithOpenApi();

            group.MapPost("/", async (Shop shop, ApplicationDbContext db) =>
            {
                var dbobj = new Shop
                {
                    OwnerUserID = shop.OwnerUserID,
                    Name = shop.Name,
                    Caption = shop.Caption
                };

                db.Shops.Add(dbobj);
                await db.SaveChangesAsync();
                return TypedResults.Created($"/api/shops/{shop.ID}", shop);
            })
            .WithName("CreateShop")
            .WithOpenApi();

            group.MapDelete("/{id:long}", async Task<Results<Ok, NotFound>> (long id, ApplicationDbContext db) =>
            {
                var dbobj = db.Shops.First(_ => _.ID == id);
                if (dbobj == null) return TypedResults.NotFound();
                dbobj.IsDeleted = true;
                await db.SaveChangesAsync();

                return TypedResults.Ok();
            })
            .WithName("DeleteShop")
            .WithOpenApi();
        }
    }
}