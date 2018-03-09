import {Injectable} from "@angular/core";
import {HttpClient, HttpHeaders} from "@angular/common/http";
import {TestSuitesComponent} from "./test-suites.component";
import {Observable} from "rxjs/Observable";
import {environment} from "../../environments/environment";

@Injectable()
export class TestSuiteService {

  private headers: HttpHeaders;
  constructor(private http: HttpClient) {
    this.headers = new HttpHeaders({"Content-Type": "application/json"});
  }

  postTestSuite$(testSuiteData, testSuitePath): Observable<any> {
    this.headers = this.headers.set("testSuitePathDir", testSuitePath);
    return this.http.post(`${environment.API_BASE_URL}test-suites`, testSuiteData, {
      headers: this.headers
    });
  }

  getTestSuite$(testSuitePath, testSuiteFileName): Observable<any> {
    this.headers = this.headers.set("testSuitePathDir", testSuitePath);
    return this.http.get(`${environment.API_BASE_URL}test-suites/${testSuiteFileName}`, {
      headers: this.headers
    });
  }

  getAllTestSuite$(testSuitePath): Observable<any> {
    this.headers = this.headers.set("testSuitePathDir", testSuitePath);
    return this.http.get(`${environment.API_BASE_URL}test-suites`, {
      headers: this.headers
    });
  }

  putTestSuite$(testSuiteData, testSuitePath, testSuiteFileName): Observable<any> {
    this.headers = this.headers.set("testSuitePathDir", testSuitePath);
    return this.http.put(
      `${environment.API_BASE_URL}test-suites/${testSuiteFileName}`,
      testSuiteData,
      {headers: this.headers}
    );
  }
}
