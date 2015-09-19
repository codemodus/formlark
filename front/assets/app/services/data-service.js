import { Http } from 'http/http';

export class DataService {

    constructor(http: Http) {
        this.http = http;
    }

    getCustomers() {
        return this.http.get('app/customers.json');
    }

}

