import { Injectable } from "@angular/core";
import { HttpModule } from "@angular/http";
import { Http } from "@angular/http";
import { HttpClient, HttpHeaders } from "@angular/common/http";
import { ConfigurationsComponent } from "./configurations/configurations.component";
import { Headers } from "@angular/http/src/headers";
import { Observable } from "rxjs/Observable";

@Injectable()
export class AutomatedUIServices {
  constructor(private http: HttpClient) {}

  private url = "http://localhost:9191/configs";

  private httpOptions = {
    headers: new HttpHeaders({
      "Content-Type": "application/json"
      // "C:/Users/a615194/go/src/github.com/xtracdev/automated-perf-test/config/"
    })
  };

  postConfig$(configData, configPath): Observable<any> {
    this.httpOptions.headers.append("configPathDir", "configPath");
    return this.http.post(this.url, configData, this.httpOptions);
  }
}
