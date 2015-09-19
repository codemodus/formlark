//import 'es6-shim';
import 'zone.js';
import 'reflect-metadata';
import { bind, FORM_BINDINGS } from 'angular2/angular2';
import { bootstrap } from 'angular2/bootstrap';
import { ROUTER_BINDINGS, LocationStrategy, HashLocationStrategy } from 'angular2/router';
import { HTTP_BINDINGS } from 'http/http';
import { AppComponent } from 'app/components/app/app-component';

bootstrap(AppComponent, [
    ROUTER_BINDINGS,
    FORM_BINDINGS,
    HTTP_BINDINGS,
    bind(LocationStrategy).toClass(HashLocationStrategy)
]);


