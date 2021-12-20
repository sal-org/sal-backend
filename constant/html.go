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
)
