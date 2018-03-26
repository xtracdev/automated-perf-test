import {Injectable} from "@angular/core";
import {HttpClient, HttpHeaders} from "@angular/common/http";
import {ConfigurationsComponent} from "./configurations.component";
import {Observable} from "rxjs/Observable";
import { environment } from '../../environments/environment';


@Injectable()
export class ConfigurationService {
  private headers: HttpHeaders;
  constructor(private http: HttpClient) {
    this.headers = new HttpHeaders({"Content-Type": "application/json"});
  }

  postConfig$(configData, configPath): Observable<any> {
    this.headers = this.headers.set("path", configPath);
    return this.http.post(`${environment.API_BASE_URL}configs`, configData, {
      headers: this.headers
    });
  }

  getConfig$(configPath, xmlFileName): Observable<any> {
    this.headers = this.headers.set("path", configPath);
    return this.http.get(`${environment.API_BASE_URL}configs/${xmlFileName}`, {
      headers: this.headers
    });
  }

  putConfig$(configData, configPath, xmlFileName): Observable<any> {
    this.headers = this.headers.set("path", configPath);
    return this.http.put(
      `${environment.API_BASE_URL}configs/${xmlFileName}`,
      configData,
      {headers: this.headers}
    );
  }

  getSchema$(location: string): Observable<any> {
    return this.http
      .get(`/${location}`, {headers: this.headers})
      .map((data: any) => data);
  }
}
