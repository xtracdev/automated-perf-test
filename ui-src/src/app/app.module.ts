import { BrowserModule } from "@angular/platform-browser";
import { NgModule } from "@angular/core";
import { BrowserAnimationsModule } from "@angular/platform-browser/animations";
import { AppComponent } from "./app.component";
import { NoopAnimationsModule } from "@angular/platform-browser/animations";
import { TestCasesComponent } from "./test-cases/test-cases.component";
import { TestSuitesComponent } from "./test-suites/test-suites.component";
import { AppRoutingModule } from "./app-routing.module";
import { NgBootstrapFormValidationModule } from "ng-bootstrap-form-validation";
import { JsonSchemaFormModule } from "angular2-json-schema-form";
import { ConfigurationsComponent } from "./configurations/configurations.component";
import { MatButtonModule, MatCheckboxModule } from "@angular/material";
@NgModule({
  declarations: [
    AppComponent,
    TestCasesComponent,
    TestSuitesComponent,
    ConfigurationsComponent
  ],
  imports: [
    JsonSchemaFormModule,
    NoopAnimationsModule,
    BrowserModule,
    NgBootstrapFormValidationModule.forRoot(),
    BrowserAnimationsModule,
    MatButtonModule,
    MatCheckboxModule,
    AppRoutingModule
  ],
  providers: [],
  bootstrap: [AppComponent]
})
export class AppModule {}
