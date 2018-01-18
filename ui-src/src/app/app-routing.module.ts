import { NgModule } from "@angular/core";
import { CommonModule } from "@angular/common";
import { RouterModule, Routes } from "@angular/router";
import { ConfigurationsComponent } from "./configurations/configurations.component";
import { TestCasesComponent } from "./test-cases/test-cases.component";
import { TestSuitesComponent } from "./test-suites/test-suites.component";

const routes: Routes = [
  { path: "", redirectTo: "/configurations", pathMatch: "full" },
  { path: "configurations", component: ConfigurationsComponent },
  { path: "test-cases", component: TestCasesComponent },
  { path: "test-suites", component: TestSuitesComponent }
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule {}
