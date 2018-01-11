import { NgModule } from "@angular/core";
import { CommonModule } from "@angular/common";

import { ConfigurationsComponent } from "./configurations.component";
import { ApplicationPropertiesComponent } from "./application-properties/application-properties.component";
import { TestCriteriaComponent } from "./test-criteria/test-criteria.component";
import { OutputPathsComponent } from "./output-paths/output-paths.component";
import { FormsModule } from "@angular/forms";
import { ReactiveFormsModule } from "@angular/forms";
@NgModule({
  imports: [CommonModule, FormsModule, ReactiveFormsModule],
  exports: [ConfigurationsComponent],
  declarations: [
    ConfigurationsComponent,
    ApplicationPropertiesComponent,
    TestCriteriaComponent,
    OutputPathsComponent
  ],

  providers: []
})
export class ConfigurationsModule {}
