import { Pipe } from 'angular2/change_detection';

export class CurrencyPipe extends Pipe {

  supports(obj) {
      return true;
  }

  transform(value) {
      if (value && !isNaN(value)) {
          return '$' + parseFloat(value).toFixed(2);
      }
      return '$0.00';
  }

  create() {
      return this;
  }
}
