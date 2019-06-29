package
{
	import flash.display.Sprite;
	import flash.text.TextField;
	import flash.text.TextFormat;
	import lesta.api.ModBase;
	import lesta.data.GameDelegate;
	import flash.external.ExternalInterface;
	
	
	/**
	 * ...
	 * @author 
	 */
	public class Main extends ModBase 
	{
		private var tf: TextField = new TextField();
		
		public function Main() {
			super();
		}
		
		override public function init(): void {
			super.init();
			
			var format: TextFormat = new TextFormat();
			format.font = 'Arial';
			format.size = 14;
			gameAPI.stage.addChild(tf);
			tf.defaultTextFormat = format;
			tf.textColor = 0xff0000;
			tf.text = "Hurr durr";
			tf.width = 800;
			tf.height = 1000;
			tf.multiline = true;
			tf.mouseEnabled = false;
			tf.selectable = false;
			
			gameAPI.data.addCallBack('sthub.LastEvent', this.onSetLastEvent);
			gameAPI.data.addCallBack('sthub.CallMeBaby', this.onInitCall);
		}
		
		override public function updateStage(width: Number, height: Number): void {
			super.updateStage(width, height);
		}
		
		private function onInitCall(): void {
			gameAPI.data.call('sthub.SomeCall', []);
		}
		
		private function onSetLastEvent(text: String): void {
			tf.text = text + '\n' + tf.text;
		}
		
	}
	
}