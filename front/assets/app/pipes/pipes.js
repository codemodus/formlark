import { defaultPipes } from 'angular2/change_detection';
import { CurrencyPipe } from './currency';

export var pipes = Object.assign({}, defaultPipes, {
  'currency': [new CurrencyPipe()]
});