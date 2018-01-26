import { BrowserModule } from "@angular/platform-browser";
import { NgModule } from "@angular/core";
import { BrowserAnimationsModule } from "@angular/platform-browser/animations";
import { AppComponent } from "./app.component";
import { NoopAnimationsModule } from "@angular/platform-browser/animations";
import { TestCasesComponent } from "./test-cases/test-cases.component";
import { TestSuitesComponent } from "./test-suites/test-suites.component";
import { AppRoutingModule } from "./app-routing.module";
import { JsonSchemaFormModule } from "angular2-json-schema-form";
import { ConfigurationsComponent } from "./configurations/configurations.component";
import { MatButtonModule, MatCheckboxModule } from "@angular/material";
import { ToastModule } from "ng2-toastr/ng2-toastr";
import { AutomatedUIServices } from "./automated-ui-services";
import { HttpClientModule, HttpClient } from "@angular/common/http";
import { FormsModule } from "@angular/forms";

@NgModule({
  declarations: [
    AppComponent,
    TestCasesComponent,
    TestSuitesComponent,
    ConfigurationsComponent
  ],
  imports: [
    FormsModule,
    JsonSchemaFormModule,
    NoopAnimationsModule,
    BrowserModule,
    BrowserAnimationsModule,
    ToastModule.forRoot(),
    MatButtonModule,
    MatCheckboxModule,
    AppRoutingModule,
    HttpClientModule
  ],
  providers: [AutomatedUIServices],
  bootstrap: [AppComponent]
})
export class AppModule {}
