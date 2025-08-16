const LOCAL_BASKET_KEY = "anonymous_basket";

interface BasketItem {
  goodId: string;
  quantity: number;
  price: number;
}

export class localBasket {
  public static getItems(): BasketItem[] {
    const items = localStorage.getItem(LOCAL_BASKET_KEY);
    return items ? (JSON.parse(items) as BasketItem[]) : ([] as BasketItem[]);
  }

  public static addItem(item: BasketItem) {
    const items = this.getItems();
    const existing = items.find((i) => i.goodId === item.goodId);
    if (existing) {
      existing.quantity += item.quantity;
    } else {
      items.push({
        goodId: item.goodId,
        quantity: item.quantity,
        price: item.price,
        //title: item.title,
        //shopTitle: item.shopTitle,
        //senderId: item.senderId,
        //shopId: item.shopId,
      });
    }
    localStorage.setItem(LOCAL_BASKET_KEY, JSON.stringify(items));
  }

  public static decItem(goodId: string) {
    let items = this.getItems();
    const idx = items.findIndex((i) => i.goodId === goodId);
    if (idx !== -1) {
      if (items[idx].quantity > 1) {
        items[idx].quantity -= 1;
      } else {
        items.splice(idx, 1);
      }
      localStorage.setItem(LOCAL_BASKET_KEY, JSON.stringify(items));
    }
  }
  public static clear() {
    localStorage.removeItem(LOCAL_BASKET_KEY);
  }
}
