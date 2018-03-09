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
import { ToastModule } from "ng2-toastr/ng2-toastr";
import { ConfigurationService } from "./configurations/configuration.service";
import { TestCaseService } from "./test-cases/test-case.service";
import { TestSuiteService } from "./test-suites/test-suite.service";
import { HttpClientModule, HttpClient } from "@angular/common/http";
import { FormsModule } from "@angular/forms";
import { TestCasesSelectionComponent } from "./shared/test-cases-selection/test-cases-selection.component";
import { TestCasesListComponent } from "./shared/test-cases-list/test-cases-list.component";

@NgModule({
  declarations: [
    AppComponent,
    TestCasesComponent,
    TestSuitesComponent,
    ConfigurationsComponent,
    TestCasesSelectionComponent,
    TestCasesListComponent
  ],
  imports: [
    FormsModule,
    JsonSchemaFormModule,
    NoopAnimationsModule,
    BrowserModule,
    BrowserAnimationsModule,
    ToastModule.forRoot(),
    AppRoutingModule,
    HttpClientModule
  ],
  providers: [ConfigurationService, TestCaseService, TestSuiteService],
  bootstrap: [AppComponent]
})
export class AppModule { }
