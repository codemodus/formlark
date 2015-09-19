import { Component, View, NgFor } from 'angular2/angular2';
import { ObservableWrapper } from 'angular2/src/facade/async';
import { Inject } from 'angular2/angular2';
import { DataService } from '../../services/data-service';
import { Sorter } from '../../utils/sorter';
import { FilterTextboxComponent } from '../filter-textbox/filter-textbox-component';
import { SortByDirective } from '../../directives/sortby/sortby-directive';

@Component({ selector: 'customers' , bindings: [DataService] })
@View({
  templateUrl: 'app/components/customers/customers-component.html',
  directives: [NgFor, FilterTextboxComponent, SortByDirective]
})
export class CustomersComponent {
  
  constructor(dataService: DataService) {
    this.title = 'Customers';
    this.filterText = 'Filter Customers:';
    this.listDisplayModeEnabled = false;
    this.displayMode = {
      Card: 0,
      List: 1
    };
    this.customers = this.filteredCustomers = [];
    
    ObservableWrapper.subscribe(dataService.getCustomers(), res => {
        this.customers = this.filteredCustomers = res.json();
    });
    
    this.sorter = new Sorter();
  }

  changeDisplayMode(displayMode) {
      switch (displayMode) {
          case this.displayMode.Card:
              this.listDisplayModeEnabled = false;
              break;
          case this.displayMode.List:
              this.listDisplayModeEnabled = true;
              break;
      }
  }

  filterChanged(data) {
    if (data) {
        data = data.toUpperCase();
        let props = ['firstName', 'lastName', 'address', 'city', 'orderTotal'];
        let filtered = this.customers.filter(item => {
            let match = false;
            for (let prop of props) {
                //console.log(item[prop] + ' ' + item[prop].toUpperCase().indexOf(data));
                if (item[prop].toString().toUpperCase().indexOf(data) > -1) {
                  match = true;
                  break;
                }
            };
            return match;
        });
        this.filteredCustomers = filtered;
    }
    else {
      this.filteredCustomers = this.customers;
    }
  }
  
  deleteCustomer(id) {
    
  }

  sort(prop) {
      this.sorter.sort(this.filteredCustomers, prop);
  }

}


