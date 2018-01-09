import { BrowserModule } from "@angular/platform-browser";
import { NgModule } from "@angular/core";
import { BrowserAnimationsModule } from "@angular/platform-browser/animations";
import { AppComponent } from "./app.component";
import { ConfigurationsComponent } from "./configurations/configurations.component";
import { TestCasesComponent } from "./test-cases/test-cases.component";
import { TestSuitesComponent } from "./test-suites/test-suites.component";
import { MatButtonModule, MatCheckboxModule } from "@angular/material";
import { AppRoutingModule } from "./app-routing.module";

import { FormsModule, ReactiveFormsModule } from "@angular/forms";
import { NgBootstrapFormValidationModule } from "ng-bootstrap-form-validation";

@NgModule({
  declarations: [
    AppComponent,
    ConfigurationsComponent,
    TestCasesComponent,
    TestSuitesComponent
  ],
  imports: [
    BrowserModule,
    FormsModule,
    ReactiveFormsModule,
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
