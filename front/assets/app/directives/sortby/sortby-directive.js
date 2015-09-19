import { Directive, ElementRef, EventEmitter } from 'angular2/angular2';

@Directive({
  selector: '[sort-by]',
  events: ['sorted']
})
export class SortByDirective {

    constructor(el: ElementRef) {
      this.sortProperty = el.nativeElement.getAttribute('sort-by');
      el.nativeElement.addEventListener('click', (event) => this.elementClicked(event));
      this.sorted = new EventEmitter();
    }

    elementClicked(event) {
        event.preventDefault();
        this.sorted.next(this.sortProperty); //Raise clicked event
    }

}
