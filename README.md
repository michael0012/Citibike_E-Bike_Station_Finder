<h1>Citibike E-Bike Station Finder</h1>

<h2>Description</h2>
<br/>
This project is meant to help people find Citibike stations with electric bikes and a limited amount of mechanical bikes. The default is set to zero mechanical bikes at the station because selecting an electric bike from a station like this has a zero additional cost for the 45 min ride as a member if you select limted assist. Currently only command line functionality exist, but looking into expanding into web app.

<h2>Usage</h2>
<br>
There are 4 possible option arguments you can use when running this program. The command line arguments are as follows:<br>
<ul>
    <li>top: This argument limits the results to the given number (default is 100)</li>
    <li>lat: This is the latitude of the area you wish to search stations near</li>
    <li>long: This is the longitude of the area you wish  to search stations near</li>
    <li>bikes: This is the number of classic or mechanical bikes you wish at most include in your search result of the station. (default is 0)</li>
</ul>
<b>Note: If lat or long are not included with the run command the default longitude and/or latitude of the search is the longitude and/or latitude of Sunset Park in Brooklyn</b><br/>
<h2>Example:</h2>
<pre>go run main.go -top 5 -lat 40.64187 -long -74.0021 -bikes 0 </pre>
This command returns the 5 closest stations with 0 classic/mechanical bikes in the starting search location of 40.64187, -74.0021
