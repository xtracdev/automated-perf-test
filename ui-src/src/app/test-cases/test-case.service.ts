import {Injectable} from "@angular/core";
import {HttpClient, HttpHeaders} from "@angular/common/http";
import {TestCasesComponent} from "./test-cases.component";
import {Observable} from "rxjs/Observable";
import { environment } from '../../environments/environment';

@Injectable()
export class TestCaseService {

  private headers: HttpHeaders;
  constructor(private http: HttpClient) {
    this.headers = new HttpHeaders({"Content-Type": "application/json"});
  }

  getTestCases$(testCasePath): Observable<any> {
    this.headers = this.headers.set("path", testCasePath);
    return this.http.get(`${environment.API_BASE_URL}test-cases`, {
      headers: this.headers
    });
  }

  postTestCase$(testCaseData, testCasePath): Observable<any> {
    this.headers = this.headers.set("path", testCasePath);
    return this.http.post(`${environment.API_BASE_URL}test-cases`, testCaseData, {
      headers: this.headers
    });
  }

  getTestCase$(testCasePath, testCaseFileName): Observable<any> {
    this.headers = this.headers.set("path", testCasePath);
    return this.http.get(`${environment.API_BASE_URL}test-cases/${testCaseFileName}`, {
      headers: this.headers
    });
  }

  putTestCase$(testCaseData, testCasePath, testCaseFileName): Observable<any> {
    this.headers = this.headers.set("path", testCasePath);
    return this.http.put(
      `${environment.API_BASE_URL}test-cases/${testCaseFileName}`,
      testCaseData,
      { headers: this.headers }
    );
  }
}
