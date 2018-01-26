import { Injectable } from "@angular/core";
import { HttpModule } from "@angular/http";
import { Http } from "@angular/http";
import { HttpClient, HttpHeaders } from "@angular/common/http";
import { ConfigurationsComponent } from "./configurations/configurations.component";
import { Headers } from "@angular/http/src/headers";
import { Observable } from "rxjs/Observable";
import { NgModel } from "@angular/forms/src/directives/ng_model";

@Injectable()
export class AutomatedUIServices {
  constructor(private http: HttpClient) {}

  private url = "http://localhost:9191/configs";
  private getUrl = "http://localhost:9191/configs/";

  postConfig$(configData, configPath): Observable<any> {
    let headers = new HttpHeaders();
    headers = headers.set("Content-Type", "application/json;");
    headers = headers.set("configPathDir", configPath);

    const httpOptions = {
      headers: headers
    };

    httpOptions.headers.append("configPathDir", configPath);
    return this.http.post(this.url, configData, httpOptions);
  }
  getConfig$(configPath, xmlFileName): Observable<any> {
    let headers = new HttpHeaders();
    headers = headers.set("Content-Type", "application/json;");
    headers = headers.set("configPathDir", configPath);

    const httpOptions = {
      headers: headers
    };
    httpOptions.headers.append("configPathDir", configPath);
    return this.http.get(this.getUrl + xmlFileName, httpOptions);
  }
}
