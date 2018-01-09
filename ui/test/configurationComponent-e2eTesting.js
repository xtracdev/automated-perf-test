describe('configuration component -e2e testing', function() {
    var goPath = os.
    var path = require('path');


     //filePath to config dir
     var configFileLocation = "/config/",
         absolutePath = path.resolve(__dirname, configFileLocation);
    //Application Properties
     var configPathDir = element(by.id("configPathDirTxt"));
     var applicationNameTxt = element(by.id('applicationNameTxt'));
     var targetHostTxt = element(by.id('targetHostTxt'));
     var targetPortTxt = element(by.id('targetPortTxt'));
     var memoryEndPointTxt = element(by.id('memoryEndPointTxt'));
     var saveBtn = element.by.id('saveBtn')
     //Test Criteria
     var numIterationsTxt = element(by.id('numIterationsTxt'));
     var concurrentUsersTxt = element(by.id('concurrentUsersTxt'));
     var memoryVarianceTxt = element(by.id('memoryVarianceTxt'));
     var serviceVarianceTxt = element(by.id('serviceVarianceTxt'));
     var testSuiteTxt = element(by.id('testSuiteTxt'));
     var requestDelayTxt = element(by.id('requestDelayTxt');
     var tpsFreqTxt = element(by.id('tpsFreqTxt');
     var rampUsersTxt = element(by.id('rampUsersTxt');
     var rampDelayTxt = element(by.id('rampDelayTxt');


     //Output paths
     var testCaseDirTxt = element(by.id('testCaseDirTxt');
     var testSuiteDirTxt = element(by.id('testSuiteDirTxt');
     var baseStatsDirTxt = element(by.id('baseStatsDirTxt');
     var reportsDirTxt = element(by.id('reportsDirTxt');


     beforeEach(function() {
        browser.get('http://localhost:9191/configs');
      });


      function addData(){

          configPathDir.sendKeys(absolutePath);
          applicationNameTxt.sendKeys("Xtrac API");
          targetHostTxt.sendKeys("localhost");
          targetPortTxt.sendKeys("9191");
          memoryEndPointTxt.sendKeys("/alt/debug/vars");

          //Populate Test Criteria
          numIterationsTxt.sendKeys(1000);
          concurrentUsersTxt.sendKeys(10);
          memoryVarianceTxt.sendKeys(15);
          serviceVarianceTxt.sendKeys(15);
          testSuiteTxt.sendKeys("suiteFileName.xml");
          requestDelayTxt.sendKeys(300);
          tpsFreqTxt.sendKeys(60);
          rampUsersTxt.sendKeys(2);
          rampDelayTxt.sendKeys(60);

          //Populate Output paths
          testCaseDirTxt.sendKeys("./definitions/testCases");
          testSuiteDirTxt.sendKeys("./definitions/testSuites");
          baseStatsDirTxt.sendKeys("./envStats");
          reportsDirTxt.sendKeys("./report");

      }


  it('should create xml file through UI', function() {
    addData();

    //Click Save
    saveBtn.click();

    browser.refresh();
    configPathDir.sendKeys(absolutePath + "Xtrac API.xml");

    //Test expected values
    expect(applicationNameTxt.getText()).toEqual('Xtrac API');
    expect(numIterationsTxt.getText()).toEqual(1000);
    expect(configPathDir.getText()).toEqual(absolutePath + "Xtrac API.xml");



  });
  
   it('should update xml file that already exists through UI', function() {
      //Populate Application Properties

     addData();
     configPathDir.sendKeys(absolutePath + "config.xml");
      //Click Save
      saveBtn.click();
  
      browser.refresh();
      configPathDir.sendKeys(absolutePath + "config.xml");
  
      //Test expected values
      expect(applicationNameTxt.getText()).toEqual('Xtrac API');
      expect(numIterationsTxt.getText()).toEqual(1100);

      expect(configPathDir.getText()).toEqual(absolutePath + "config.xml");
  
  
  
    });

  it('should not enable save button if all req fields are not filled (applicationName) ', function() {
    addData();
    applicationNameTxt.sendKeys(null);


    //Test
    expect(saveBtn.isEnabled()).toBe(false);


  });

    it('should not enable save button if all req fields are not filled (targetPort) ', function() {

      addData();
      targetPortTxt.sendKeys(null);


      //Test
      expect(saveBtn.isEnabled()).toBe(false);


    });
        it('should not enable save button if all req fields are not filled (targetHost) ', function() {
            addData();
          targetHostTxt.sendKeys(null);

          //Test
          expect(saveBtn.isEnabled()).toBe(false);


        });




    it('should not enable save button if all req fields are not filled (configPathDir) ', function() {

      addData();
      configPathDir.sendKeys(null);


      //Test
      expect(saveBtn.isEnabled()).toBe(false);





    });

        it('should not enable save button if all req fields are not filled (numIterations) ', function() {
          addData
          numIterationsTxt.sendKeys(null);


          //Test
          expect(saveBtn.isEnabled()).toBe(false);


        });

        it('should not enable save button if all req fields are not filled (concurrentUsers) ', function() {
                   addData();
                  numIterationsTxt.sendKeys(1000);


                  //Test
                  expect(saveBtn.isEnabled()).toBe(false);


           });

        it('should not enable save button if all req fields are not filled (memoryVariance) ', function() {
                  addData();
                  memoryVarianceTxt.sendKeys(null);


                  //Test
                  expect(saveBtn.isEnabled()).toBe(false);


           });
        it('should not enable save button if all req fields are not filled (serviceVariance) ', function() {
                   addData();
                  serviceVarianceTxt.sendKeys(null);


                  //Test
                  expect(saveBtn.isEnabled()).toBe(false);


               });

     it('should not enable save button if all req fields are not filled (testSuite) ', function() {
                  addData();
                  testSuiteTxt.sendKeys(null);


                  //Test
                  expect(saveBtn.isEnabled()).toBe(false);


               });

     it('should not enable save button if all req fields are not filled (requestDelay) ', function() {
                      addData();
                      requestDelayTxt.sendKeys(null);


                      //Test
                      expect(saveBtn.isEnabled()).toBe(false);


                   });



     it('should not enable save button if all req fields are not filled (tpsFreq) ', function() {
                      addData();
                      tpsFreqTxt.sendKeys(null);


                      //Test
                      expect(saveBtn.isEnabled()).toBe(false);


                   });


        it('should not enable save button if all req fields are not filled (tpsFreq) ', function() {
                       addData();
                      tpsFreqTxt.sendKeys(null);


                      //Test
                      expect(saveBtn.isEnabled()).toBe(false);


                   });


              it('should not enable save button if all req fields are not filled (rampUsers) ', function() {
                                  //Populate Application Properties
                                  configPathDir.sendKeys(absolutePath);
                                  applicationNameTxt.sendKeys("Test");
                                  targetPortTxt.sendKeys("localhost");
                                  memoryEndPointTxt.sendKeys("/alt/debug/vars");
                                  targetHostTxt.sendKeys("localhost");

                                  //Populate Test Criteria
                                  numIterationsTxt.sendKeys(1000);
                                  concurrentUsersTxt.sendKeys(10);
                                  memoryVarianceTxt.sendKeys(15);
                                  serviceVarianceTxt.sendKeys(15);
                                  testSuiteTxt.sendKeys("suiteFileName.xml");
                                  requestDelayTxt.sendKeys(300);
                                  tpsFreqTxt.sendKeys(5);
                                  rampUsersTxt.sendKeys(null);
                                  rampDelayTxt.sendKeys(60);

                                  //Populate Output paths
                                  testCaseDirTxt.sendKeys("./definitions/testCases");
                                  testSuiteDirTxt.sendKeys("./definitions/testSuites");
                                  baseStatsDirTxt.sendKeys("./envStats");
                                  reportsDirTxt.sendKeys("./report");

                                  //Test
                                  expect(saveBtn.isEnabled()).toBe(false);


                               });


              it('should not enable save button if all req fields are not filled (rampDelay) ', function() {
                                  addData();
                                  rampDelayTxt.sendKeys(null);

                                  //Test
                                  expect(saveBtn.isEnabled()).toBe(false);


                               });



        it('should not enable save button if all req fields are not filled (testCaseDir) ', function() {
                                          addData();

                                       //Populate Output paths
                                       testCaseDirTxt.sendKeys(null);


                                       //Test
                                       expect(saveBtn.isEnabled()).toBe(false);


                                    });


            it('should not enable save button if all req fields are not filled (testSuiteDirTxt) ', function() {
                                      addData();
                                      testSuiteDirTxt.sendKeys(null);


                                      //Test
                                      expect(saveBtn.isEnabled()).toBe(false);


                                   });



          it('should not enable save button if all req fields are not filled (baseStatsDir) ', function() {
                                         addData();
                                         baseStatsDirTxt.sendKeys(null);


                                         //Test
                                         expect(saveBtn.isEnabled()).toBe(false);


                                      });



     it('should not enable save button if all req fields are not filled (reportsDirTxt) ', function() {
                                      addData();

                                      reportsDirTxt.sendKeys(null);

                                      //Test
                                      expect(saveBtn.isEnabled()).toBe(false);

                            });
   it('should create file when memoryEndppoint field is empty ', function() {

        addData();
      //Left Blank
       memoryEndPointTxt.sendKeys(null);




       //Click Save
        saveBtn.click();

        browser.refresh();
        configPathDir.sendKeys(absolutePath + "Xtrac API.xml");

        //Test expected values
        expect(applicationNameTxt.getText()).toEqual('Xtrac API');
        expect(numIterationsTxt.getText()).toEqual(1000);
        expect(configPathDir.getText()).toEqual(absolutePath + "XtracTest.xml");


   });

   });