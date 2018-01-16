import { Injectable } from "@angular/core";
import { HttpModule } from "@angular/http";
import { Http } from "@angular/http";
import { HttpClient, HttpHeaders } from "@angular/common/http";
import { ConfigurationsComponent } from "./configurations/configurations.component";
import { Headers } from "@angular/http/src/headers";

const httpOptions = {
  headers: new HttpHeaders({
    configPathDir:
      "C:/Users/a615194/go/src/github.com/xtracdev/automated-perf-test/config/"
  })
};

@Injectable()
export class PostService {
  constructor(private http: HttpClient) {}

  private url = "http://localhost:9191/configs";

  addConfig(form: FormData): void {
    console.log(form);
    this.http.post<FormData>(this.url, form, httpOptions);
  }
}
