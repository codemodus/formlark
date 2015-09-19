import { Component, View, LifecycleEvent, EventEmitter } from 'angular2/angular2';
import { FORM_DIRECTIVES } from 'angular2/forms';


@Component({
  selector: 'filter-textbox',
  events: ['changed'],
  properties: ['text'],
  lifecycle: [LifecycleEvent.onChange]
})
@View({
  template: `
    <form>
         Filter:
         <input type="text" 
                [(ng-model)]="model.filter" 
                (keyup)="filterChanged($event)"  />
    </form>
  `,
  directives: [FORM_DIRECTIVES]
})
export class FilterTextboxComponent {

    constructor() {
      this.model = {
        filter: null
      };
      this.changed = new EventEmitter();
    }

    filterChanged(event) {
        event.preventDefault();
        this.changed.next(this.model.filter); //Raise changed event
    }

    onChange(changes) {
      //alert(changes);
    }

}
