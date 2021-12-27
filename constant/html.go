package constant

const (
	CounsellorProfileHtml = `<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<meta http-equiv="X-UA-Compatible" content="IE=edge,chrome=1">
	</head>
	<body>
		<h3 style="color: rgb(113, 113, 233);">Counsellor Profile Details</h3>
		<table style="width: 50%;text-align: left;">
		  <tr>
			<th style="color:#808080; padding: 10px;">Counsellor Name:</th>
			<td>###First_Name###</td><td>###Last_Name###</td>
		  </tr>
			<tr>
				<th style="color:#808080; padding: 10px;">Gender:</th>
				<td>###Gender###</td>
			  </tr>
			  <tr>
				<th style="color:#808080; padding: 10px;">Phone Number:</th>
				<td>###Phone###</td>
			  </tr>
			  <tr>
				<th style="color:#808080; padding: 10px;">Photo:</th>
				<td><img src="https://sal-prod.s3.ap-south-1.amazonaws.com/###Photo###"></td>
			  </tr>
			  <tr>
				<th style="color:#808080; padding: 10px;">Email:</th>
				<td>###Email###</td>
			  </tr>
			  <tr>
				<th style="color:#808080; padding: 10px;">Education:</th>
				<td>###Education###</td>
			  </tr>
			  <tr>
				<th style="color:#808080; padding: 10px;">Experience:</th>
				<td>###Experience###</td>
			  </tr>
			  <tr>
				<th style="color:#808080; padding: 10px;">About:</th>
				<td>###About###</td>
			  </tr>
			  <tr>
				<th style="color:#808080; padding: 10px;">Resume:</th>
				<td><a href="https://sal-prod.s3.ap-south-1.amazonaws.com/###Resume###">Resume Document</a></td>
			  </tr>
			  <tr>
				<th style="color:#808080; padding: 10px;">Certificate:</th>
				<td><a href="https://sal-prod.s3.ap-south-1.amazonaws.com/###Certificate###">Certificate Document</a></td>
			  </tr>
			  <tr>
				<th style="color:#808080; padding: 10px;">Aadhar:</th>
				<td><a href="https://sal-prod.s3.ap-south-1.amazonaws.com/###Aadhar###">Aadhar Document</a></td>
			  </tr>
			  <tr>
				<th style="color:#808080; padding: 10px;">Linkedin:</th>
				<td><a href="###Linkedin###">Linkedin Profile</a></td>
			  </tr>
			  <tr>
				<th style="color:#808080; padding: 10px;">Status:</th>
				<td>###Status###</td>
			  </tr>
		</table>
	</body>
	</html>`

	EventWaitingForApprovalBody = `<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<meta http-equiv="X-UA-Compatible" content="IE=edge,chrome=1">
	</head>
	<body>
		<h3 style="color: rgb(113, 113, 233);">SAL Cafe Details</h3>
		<table style="width: 50%;text-align: left;">
		  <tr>
			<th style="color:#808080; padding: 10px;">Counsellor Name:</th>
			<td>###First_name###</td><td>###Last_name###</td>
		  </tr>
			<tr>
				<th style="color:#808080; padding: 10px;">Type:</th>
				<td>###Type###</td>
			  </tr>
			  <tr>
				<th style="color:#808080; padding: 10px;">SAL Cafe Name:</th>
				<td>###Title###</td>
			  </tr>
			  <tr>
				<th style="color:#808080; padding: 10px;">Event Description:</th>
				<td>###Description###</td>
			  </tr>
			  <tr>
				<th style="color:#808080; padding: 10px;">Event Photo:</th>
				<td><img src="https://sal-prod.s3.ap-south-1.amazonaws.com/###Photo###"></td>
			  </tr>
			  <tr>
				<th style="color:#808080; padding: 10px;">Event Topic:</th>
				<td>###Topic_id###</td>
			  </tr>
			  <tr>
				<th style="color:#808080; padding: 10px;">Event Date:</th>
				<td>###Date###</td>
			  </tr>
			  <tr>
				<th style="color:#808080; padding: 10px;">Event Time:</th>
				<td>###Time###</td>
			  </tr>
			  <tr>
				<th style="color:#808080; padding: 10px;">Event Duration:</th>
				<td>###Duration###</td>
			  </tr>
			  <tr>
				<th style="color:#808080; padding: 10px;">Event Price:</th>
				<td>###Price###</td>
			  </tr>
		</table>
	</body>
	</html>`

	SendAReceiptForClient = `<!DOCTYPE html>
	<html>
	<head>
		<meta charset='utf-8'>
		<meta http-equiv='X-UA-Compatible' content='IE=edge'>
		<title>Receipt Page</title>
		<meta name='viewport' content='width=device-width, initial-scale=1'>
		
		<style>
			body{
				min-width:300px;        /* Suppose you want minimum width of 1000px */
				width: auto !important;  /* Firefox will set width as auto */
				width:300px;            /* As IE6 ignores !important it will set width as 1000px; */
				
			}
	
			th{
				color:#38373781;
				text-align: left;
			}
			td{
				text-align:right ;
				color:#38373781;
			}
			@media screen and (max-width:355px) {
				.footer-items{
					flex-direction: column;
					align-items: center;
					justify-content: center;
				}
			}
			
	
		</style>
	</head>
	<body>
	  <div id="maindiv;" style="padding:2%">
	
		<div>
			<h2 style="text-align:left;font-size:230%;">
				Receipt
				<span style="float: right; color:#0066B3;font-size:180% ;font-weight: 100;">
					SAL
				</span>
			</h2>
		</div>
		
		<br>
		<hr style="background-color: #0066B3;height: 1px;border: 0px;"/>
	
		<div>
			<h3 style="color: #333333;text-align: center;">
				Thank you for engaging with SAL for your Emotional Wellbeing.
			</h3>
			<h4 style="color: #333333;">
				Receipt Details
			</h4>
	
			<table style="width:100% ; text-align: left;">
				<tr>
				<th style=" padding: 10px;">Date:</th>
				<td>###Date###</td>
				</tr>
				<tr>
				<th style="padding: 10px;">Receipt No:</th>
				<td>###ReceiptNo###</td>
				</tr>
				<tr>
				<th style="padding: 10px;">Reference No:</th>
				<td>###ReferenceNo###</td>
				</tr>
			</table>
		</div>
	
		<br/>
		<hr style="background-color: #0066B3;height: 1px;border: 0px;"/>
		
		<div>
			<h4>
				Summary
			</h4>
			<table style="width:60%; float: right;">
				<tr>
				<th style="padding: 10px;text-align: right; ">Price</th>
				<th style="padding: 10px;text-align: right;">Qty</th>
				<th style="padding: 10px;text-align: right;">Total</th>
				</tr>
				<tr>
				<td style="padding-top: 20px;text-align: center;">###SPrice###</td>
				<td style="padding-top: 20px;text-align: center;">###Qty###</td>
				<td style="padding-top: 20px;text-align: center;">###Total###</td>
				</tr>
			</table>
			<p style="padding-top: 40px;margin-top: 40px;color:#38373781;">###SessionsType###</p>
		</div>
	
		<br/>
		<hr style="background-color: #0066B3;height: 1px;border: 0px;"/>
	
		<div>
			<table style="width:100% ;">
				<tr >
				<th style="text-align: left;color: #000;"> Transaction Details</th>
				<th style="text-align: right;color: #000;">INR</td>
				</tr>
				<tr >
					<th style="padding-top: 10px;">Price:</th>
				  <td style="padding-top: 10px;">###TPrice###</td>
				</tr>
				<tr>
				  <th style="padding-top: 10px;">Discount:</th>
				  <td style="padding-top: 10px;color: rgb(245, 79, 79);">- ###Discount###</td>
				</tr>
				<tr style="border: 4px ridge">
				  <th style="color:#0066B3;padding-top: 10px;">Total Paid <br>(All Inclusive):</th>
				  <td style="color:#0066B3;padding-top:10px; font-weight: bold;">Rs. ###TotalP###</td>
				</tr>
			</table>
		</div>
	
		<br>
		<hr style="background-color: #0066B3;height: 1px;border: 0px;"/>
	
		<div class="footer-items" style="display: flex;align-items: flex-start;justify-content: space-between;">
			<div >
				<h4>Salubrium Private Limited</h4>
				<h5><a href="www.sal-foundation.com">www.sal-&nbsp;foundation.com</a></h5>
				<h5>reachus@sal.foundation</h5>
			</div>
			<div >
				<h5>Follow us at<h5>
				<img src="https://sal-prod.s3.ap-south-1.amazonaws.com/content/SMicons.jpg" alt="Social mediun Logo"/>
			</div>
			<div >
			  <img src="https://sal-prod.s3.ap-south-1.amazonaws.com/content/sal-tagline.png" alt="SAL Tag" style="object-fit: cover;width: 80%;height:80%;"/>
			</div>
		</div>
	
	  </div>  
	</body>
	</html>`
)
