import { BrowserModule } from "@angular/platform-browser";
import { NgModule } from "@angular/core";
import { BrowserAnimationsModule } from "@angular/platform-browser/animations";
import { AppComponent } from "./app.component";
import { NoopAnimationsModule } from "@angular/platform-browser/animations";
import { TestCasesComponent } from "./test-cases/test-cases.component";
import { TestSuitesComponent } from "./test-suites/test-suites.component";
import { MatButtonModule, MatCheckboxModule } from "@angular/material";
import { AppRoutingModule } from "./app-routing.module";
import { FormsModule, ReactiveFormsModule } from "@angular/forms";
import { NgBootstrapFormValidationModule } from "ng-bootstrap-form-validation";
import { ConfigurationsModule } from "./configurations/configurations.module";
import { JsonSchemaFormModule } from "angular2-json-schema-form";

@NgModule({
  declarations: [AppComponent, TestCasesComponent, TestSuitesComponent],
  imports: [
    JsonSchemaFormModule,
    NoopAnimationsModule,
    BrowserModule,
    NgBootstrapFormValidationModule.forRoot(),
    BrowserAnimationsModule,
    MatButtonModule,
    MatCheckboxModule,
    AppRoutingModule,
    ConfigurationsModule
  ],
  providers: [],
  bootstrap: [AppComponent]
})
export class AppModule {}
