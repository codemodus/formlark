import { Component, View, NgFor } from 'angular2/angular2';
import { RouterLink } from 'angular2/router';
import { DataService } from 'app/services/data-service';

@Component({ selector: 'orders' })
@View({
  templateUrl: 'app/components/orders/orders-component.html',
  directives: [NgFor, RouterLink]
})
export class OrdersComponent {
    constructor(dataService: DataService) {
      this.title = 'Orders';
      //this.orders = dataService.getOrders();
    }
}
