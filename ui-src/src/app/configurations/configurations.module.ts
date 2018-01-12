import { NgModule } from "@angular/core";
import { CommonModule } from "@angular/common";
import { HttpModule } from "@angular/http";
import { ConfigurationsComponent } from "./configurations.component";

@NgModule({
  imports: [CommonModule, HttpModule],
  exports: [ConfigurationsComponent],
  declarations: [ConfigurationsComponent],

  providers: []
})
export class ConfigurationsModule {}
