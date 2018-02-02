import { Injectable } from "@angular/core";
import { HttpModule } from "@angular/http";
import { Http } from "@angular/http";
import { HttpClient, HttpHeaders } from "@angular/common/http";
import { ConfigurationsComponent } from "./configurations/configurations.component";
import { Headers } from "@angular/http/src/headers";
import { Observable } from "rxjs/Observable";
import { environment } from "../environments/environment.prod";
import { $ } from "protractor";


@Injectable()
export class AutomatedUIServices {

  constructor(private http: HttpClient) {}


  postConfig$(configData, configPath): Observable<any> {
    let headers = new HttpHeaders();
    headers = headers.set("Content-Type", "application/json;");
    headers = headers.append("configPathDir", configPath);
    return this.http.post(`${environment.apiBaseUrl}configs`,configData, {headers});

  }
  getConfig$(configPath, xmlFileName): Observable<any> {
    let headers = new HttpHeaders();
    headers = headers.set("Content-Type", "application/json;");
    headers = headers.append("configPathDir", configPath);
    return this.http.get(`${environment.apiBaseUrl}configs/${xmlFileName}`, {headers});
  }
  putConfig$(configData, configPath, xmlFileName): Observable<any> {
    let headers = new HttpHeaders();
    headers = headers.set("Content-Type", "application/json;");
    headers = headers.append("configPathDir", configPath);
    return this.http.put(`${environment.apiBaseUrl}configs/${xmlFileName}`, configData, {headers});

  }
}

